require 'thread'

class Listing
    require 'net/http'
    require 'json'
    require 'set'
    @@pastebin_listing_url = URI("https://scrape.pastebin.com/api_scraping.php")
    @@pastebin_listing_params = URI.encode_www_form({limit:"30"})

    public
    def initialize
        @url = @@pastebin_listing_url
        @url.query = @@pastebin_listing_params
        @listing = Set.new
    end

    private
    def get_new_listings 
        ret = nil
        begin
            response = HTTP.get_response(@url)
            tags = Set.new
            JSON.parse(@response).each {|x| tags.merge x["key"]} if response.is_a? HTTPSuccess
            ret = tags - @listing
            @listing = tags
        rescue ParserError
            puts "Error while trying to parse json"
        rescue OpenTimeout
            puts "Error timed out during request"
        end
        return ret
    end
    
end

class Scraper
    @@pastebin_scrape_url = URL("https://scrape.pastebin.com/api_scrape_item.php")
    @@pastebin_scrape_params = {i: @listing_id}

    def initialize listing_id
        @url = @@pastebin_scrape_url
        @url.query = @@pastebin_scrape_params
        @listing_id = listing_id
    end

    public
    def get_paste
        @contents = HTTP.get_response(@url)
    end

    public 
    def 
end

class Send
    def initialize title message
        @title = title
        @message = message
    end

    public
    def send
        post_paste
    end
end

class Email < send
    require 'net/smtp'

    def initialize title, message, server, src_email, dst_email, password
        super title, message
        @password = password
        @server = server
        @src_email = src_email
        @dst_email = dst_email

        @email = <<END_OF_MESSAGE
FROM: #{@src_email} <#{@src_email}>
TO: listeners <#{@dst_email}>
SUBJECT: #{@title}
DATE: #{TIME.now}

#{@message}

END_OF_MESSAGE
    end

    private 
    def post_paste
        SMTP.enable_starttls.start(@server, 587) do |smtp|
            smtp.send_message @email, @src_email, @dst_email
        end
    end
end

def get_and_send id, con
    message = Scraper.new(x).get_paste.filter
    Email.new(id, message, con["server"], con["src_email"], con["dst_email"], con["password"]).send if message.is_a? HTTPSuccess
end

def main
    #parse args here
    connection_info = {
        server: nil,
        src_email: nil
        dst_email: nil
        password: nil
    }
    pastes = Listing.new
    
    loop do
        threads = []
        pastes.get_new_listings.each {|x| threads << Thread.new {get_and_send(x, connection_info)} }
        sleep(10)
        threads.each {|x| x.join()}
    end
end