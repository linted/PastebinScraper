#!/usr/bin/ruby -w
require "url_types"

#Abstract Base Classes
class Match
    def initialize title, contents
        @title = title
        @contents = contents
        @title = ''
        @comments = {}
    end

    def find
        @title << check_title << check_contents if not blacklisted?
    end

    def check_title
        ""        #we don't care what the title says
    end

    def check_contents
        ""
    end

    def blacklisted?
        false
    end

end

#########################################
# Implementation
#########################################

class Match_url < Match
    @@url_classifier = /\b((?:https?|ftp|file):\/\/[\da-z\.-]+\.[a-z\.]{2,6}[\/\w \.-]*\/?)\b/
    @@url_types = {
        "Pastebin_Url" => URL_pastebin,
        "Imgur_Url" => URL_imgur,
        "Google_Drive" => URL_google_drive,
        "Short_link" => URL_short_link
    }

    def check_contents
        found = false
        @contents.scan(@@url_classifier) do |x|
            found = true
            check_url_type x
        end
        found ? "url" : ""
    end

    def check_url_type url

    end

end


class Match_email < Match
    @@email_classifier = /\b((([!#$%&'*+\-\/=?^`{|}~\w])|([!#$%&'*+\-\/=?^`{|}~\w][!#$%&'*+\-\/=?^`{|}~\.\w]{0,}[!#$%&'*+\-\/=?^`{|}~\w]))[@]\w+([-.]\w+)*\.\w+([-.]\w+)*)\b/

    def check_contents

    end
end
