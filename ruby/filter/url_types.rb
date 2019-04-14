
################
## ABC
################
class Sub_Match
    def initialize
        @title = ''
        @comments = ''
    end

    require 'net/http'
    def matches? contents
        @title = @@regex.match contents
        return get_parsed_matches contents
    end

    def get_parsed_matches contents 
        #This is meant to be overwritten with additional parsing
        #useful if you want to follow the link or something
        #should return True or False
        @title != ''
    end
end

###################
## Implementations
###################

class URL_short_link < Sub_Match
    @@regex = /\b(https?:\/\/
        (?:bit\.ly
        | t\.co
        | lnkd\.in
        | tcrn\.ch
        | ecleneue\.com
        | swarife\.com
        | goo\.gl
        )\S*)\b/x
    def get_parsed_matches contents
        if @title != ''
            c = []
            contents.scan(@@regex) do
                response = HTTP.get_response(URI(url))
                if response.is_a? Net::HTTPRedirection
                    c << response["location"]
                end
            end
            @comments = c.join
        else
            return false
        end
        return true
    end

end

class URL_pastebin
    @@regex = /pastebin.com/i
end

class URL_imgur
    @@regex = /imgur.com/i
end

class URL_google_drive
    @@regex = /drive.google.com/i
end