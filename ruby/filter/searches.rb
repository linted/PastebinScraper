#!/usr/bin/ruby -w
require "url_types"

#Abstract Base Classes
class Match
    def initialize title, contents
        @title = title
        @contents = contents
        @matches = ''
    end

    def find
        @matches << check_title << check_contents if not blacklisted?
    end

    def check_title
        ""        #we don't care what the title says
    end

    def check_contents
        ""
    end

end

#########################################
# Implementation
#########################################

class Match_url < Match
    @@url_classifier = /\b((https?|ftp|file):\/\/)([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?\b/
    @@url_types = {
        "Pastebin_Url" => URL_pastebin,
        "Imgur_Url" => URL_imgur,
        "Google_Drive" => URL_google_drive,
        "Short_link" => URL_short_link
    }

    def check_contents
        @@url_types.each do |title, search|
            
        end
        @matches
    end

end

