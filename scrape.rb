#!/usr/bin/ruby -w
require 'thread'
require 'optparse'
require 'io/console'

class Listing
    require 'net/http'
    require 'json'
    require 'set'
    @@pastebin_listing_url = URI("https://scrape.pastebin.com/api_scraping.php?limit=30")

    public
    def initialize
        @listing = Set.new
    end

    public
    def get_new_listings 
        ret = nil
        begin
            response = Net::HTTP.get_response(@@pastebin_listing_url)
            tags = Set.new
            
            JSON.parse(response.body ).each do |x|
                tags.add x["key"]
            end if response.is_a? Net::HTTPSuccess
            ret = tags - @listing
            @listing = tags
        rescue JSON::ParserError
            sprint {puts "Error while trying to parse json"}
        rescue Net::OpenTimeout
            sprint {puts "Error timed out during request" }
        end
        return ret
    end
    
end

class Scraper
    require "net/http"
    @@pastebin_scrape_url = "https://scrape.pastebin.com/api_scrape_item.php"
    @@searches = {
        "Email_Address" => /\b((([!#$%&'*+\-\/=?^`{|}~\w])|([!#$%&'*+\-\/=?^`{|}~\w][!#$%&'*+\-\/=?^`{|}~\.\w]{0,}[!#$%&'*+\-\/=?^`{|}~\w]))[@]\w+([-.]\w+)*\.\w+([-.]\w+)*)\b/,
        "IP_Address" => /\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b/,
        "Phone_Number" => /\b\(\d{3}\) ?\d{3}( |-)?\d{4}|^\d{3}( |-)?\d{3}( |-)?\d{4}\b/,
        "URL" => /\b((https?|ftp|file):\/\/)([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?\b/,
        "Credit Card" => /\b
                (?:4[0-9]{12}(?:[0-9]{3})?          # Visa
                |  (?:5[1-5][0-9]{2}                # MasterCard
                    | 222[1-9]|22[3-9][0-9]|2[3-6][0-9]{2}|27[01][0-9]|2720)[0-9]{12}
                |  3[47][0-9]{13}                   # American Express
                |  3(?:0[0-5]|[68][0-9])[0-9]{11}   # Diners Club
                |  6(?:011|5[0-9]{2})[0-9]{12}      # Discover
                |  (?:2131|1800|35\d{3})\d{11}      # JCB
            )\b/x
    }

    attr_reader :contents
    attr_reader :matches

    def initialize listing_id
        @url = URI(@@pastebin_scrape_url)
        @url.query = URI.encode_www_form({:i => listing_id})
        @listing_id = listing_id
        @contents = nil
        @matches = ""
    end

    public
    def get_paste
        response = Net::HTTP.get_response(@url)
        @contents = response.body if response.is_a? Net::HTTPSuccess
        self
    end

    public
    def filter
        @matches = ""
        @@searches.each {|type, pattern| @matches << type << " " if pattern.match(@contents) }
        self
    end
end

class Send
    def initialize id, title, message
        @id = id
        @title = title
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

    def initialize id, title, message, server, src_email, dst_email, password
        super id, title, message
        @password = password
        @server = server
        @src_email = src_email
        @dst_email = dst_email

        setup unless @@smtp

        @email = <<END_OF_MESSAGE
FROM: Pastebin Scraper <#{@src_email}>
TO: listeners <#{@dst_email}>
SUBJECT: #{@title}
DATE: #{Time.now}

link: https://pastebin.com/#{@id}

#{@message}

END_OF_MESSAGE
    end

    private 
    def post_paste
        @@mutex.synchronize {
            connect unless @@connection
            loop do
                begin 
                    @@connection.send_message @email, @src_email, @dst_email
                rescue StandardError => e
                    sprint { puts "Caught exception while trying to send email: #{e}"}
                    connect
                    break
                else
                    break #no errors!
                end
            end
        }
    end

    private
    def connect
        @@connection = @@smtp.start(@server, @src_email, @password, :login)
    end

    private
    def reconnect
        @@connection.finish
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

def get_and_send id, con
    sprint {puts "Starting #{id}"}
    message = Scraper.new(id).get_paste.filter
    Email.new(id, message.matches, message.contents, con[:server], con[:src_email], con[:dst_email], con[:password]).send if message.matches != ''
    sprint {puts "Finished #{id}"}
    return
end

def sprint
    /#{$mutex = Mutex.new}/o
    $mutex.synchronize {
        yield 
    }
end

def main
    #parse args here
    connection_info = {
        server: nil,
        src_email: nil,
        dst_email: nil,
        password: nil
    }
    OptionParser.new do |opts|
        opts.on("-eEMAIL", "--send-email EMAIL", "Email to send from", :REQUIRED) do |x| connection_info[:src_email] = x end
        opts.on("-rEMAIL", "--recv-email EMAIL", "Email to send to", :REQUIRED) do |x| connection_info[:dst_email] = x end
        opts.on("-sSERVER", "--smtp-server SERVER", "Smtp server to talk to", :REQUIRED) do |x| connection_info[:server] = x end
    end.parse!
    
    print "Password: "
    connection_info[:password] = STDIN.noecho(&:gets).chomp
    puts
    connection_info.each {|k,v| raise "Error, please supply all paramaters" if not v}

    pastes = Listing.new
    
    begin
        loop do
            new_pastes = pastes.get_new_listings
            sprint {puts "#{new_pastes.length} New; #{Thread.list.length - 1} running"}
            new_pastes.each {|x| Thread.new {get_and_send(x, connection_info)} }
            sleep(10)
            Thread.list.each {|x| x.join if not x.alive?} #clean up, clean up, everyone, everywhere
        end
    rescue Interrupt
        sprint {puts "Caught exception. Shutting down #{Thread.list.length - 1} threads cleanly"}
        Thread.list.each {|x| x.join unless x == Thread.current}
        sprint {puts "Threads remaining #{Thread.list.length - 1}"}
    end
end

main