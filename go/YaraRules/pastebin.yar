/*
    pastebin links
*/

rule pastebin_url
{
    meta:
        author = "linted"
        info = "pastebin url"

    strings:
        $url = /pastebin.com/
    condition:
        $url
}