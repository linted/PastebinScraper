#!/usr/bin/ruby -w

class Send
    def initialize title, id, subject, message
        @id = id
        @title = title
        @subject = subject
        @message = message
    end

    public
    def send
        post_paste
    end
end

class Email < Send
    require 'net/smtp'
    @@mutex = Mutex.new
    @@smtp = nil
    @@connection = nil

    def initialize title, id, subject, message, server, src_email, dst_email, password
        super title, id, subject, message
        @password = password
        @server = server
        @src_email = src_email
        @dst_email = dst_email

        setup unless @@smtp

        begin
            @email = <<END_OF_MESSAGE
FROM: Pastebin Scraper <#{@src_email}>
TO: listeners <#{@dst_email}>
SUBJECT: [#{@subject}] #{@title} 
DATE: #{Time.now}

link: https://pastebin.com/#{@id}

#{@message.force_encoding("UTF-8")}

END_OF_MESSAGE
        rescue Encoding::CompatibilityError => e
            sprint {puts "Error [#{id}]: #{e}"}
        end

    end

    private 
    def post_paste
        @@mutex.synchronize {
            connect unless @@connection
            loop do
                begin 
                    @@connection.send_message @email, @src_email, @dst_email
                rescue Net::SMTPUnknownError => e
                    sprint { puts "Unkown error occured: #{e.message}"}
                    $GLOBAL_STOP_FLAG = true

                rescue StandardError => e
                    sprint { puts "Error during send [#{e.class}]: #{e.message}"}
                    reconnect
                    break
                else
                    break #no errors!
                end
            end
        } unless $GLOBAL_STOP_FLAG
    end

    private
    def connect
        begin
            @@connection = @@smtp.start(@server, @src_email, @password, :login)
        rescue Net::SMTPAuthenticationError => e
            sprint {puts "Fatal Error: #{e}"}
            $GLOBAL_STOP_FLAG = true
            exit
        end            
    end

    private
    def reconnect
        begin
            @@connection.finish
        rescue StandardError
            setup
        end
        connect
    end

    public
    def setup
        @@smtp = Net::SMTP.new(@server,587)
        @@smtp.enable_starttls
    end

    public
    def shutdown
        @@connection.finish if @@connection
    end

end