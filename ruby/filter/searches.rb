#!/usr/bin/ruby -w


#Abstract Base Classes
class Match
    def initialize title, contents
        @title = title
        @contents = contents
    end

    def found?
        check_title or check_contents
    end

end

class Filter
    def initialize title, contents
        @title = title
        @contents = contents
    end

    def blacklisted?
        check_title and check_contents
    end

end

#########################################
# Implementation
#########################################

class Match_url < Match
    @@url_classifier = /\b((https?|ftp|file):\/\/)([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?\b/
    @@url_static = {
        "Pastebin_Url" => /pastebin.com/i,
        "Imgur_Url" => /imgur.com/i,
        "Google_Drive" => /drive.google.com/i
    }
    @@url_shortened = /\b(http:\/\/(?:bit\.ly|t\.co|lnkd\.in|tcrn\.ch)\S*)\b/

    def check_title
        true        #we don't care what the title says
    end

    def check_contents

    end


end