

class listing
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
            listing = JSON.parse(@response) if response.is_a? HTTPSuccess
            tags = Set.new
            listing.each {|x| tags.merge x['key']} if listing.is_a? Array
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

class scraper
    @@pastebin_scrape_url = URL("https://scrape.pastebin.com/api_scrape_item.php")
    @@pastebin_scrape_params = {i: @listing_id}

    def initialize listing_id
        @url = @@pastebin_scrape_url
        @url.query = @@pastebin_scrape_params
        @listing_id = listing_id
    end

    public
    def get_paste
        response = 


end

class send
    def send
        post_paste
    end
end

class email < send
    def initialize
    end

    private 
    def post_paste

    end
end