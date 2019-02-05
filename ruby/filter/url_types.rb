
################
## ABC
################
class Sub_Match
    require 'net/http'
    def matches? string
        @matches = @@regex.match string
        return get_parsed_matches
    end

    def get_parsed_matches
        #This is meant to be overwritten with additional parsing
        #useful if you want to follow the link or something
        @matches
    end
end

###################
## Implementations
###################

class URL_short_link < Sub_Match
    @@regex = /\b(https?:\/\/(?:bit\.ly|t\.co|lnkd\.in|tcrn\.ch|ecleneue\.com)\S*)\b/
    def get_parsed_matches
        response = HTTP.get_response(URI(url))
        if response.is_a? Net::HTTPRedirection
            return response["location"]
        end
        return @matches
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