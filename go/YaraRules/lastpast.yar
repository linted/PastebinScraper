/*
    Weird lastpast thing
*/

rule lastpast
{
    meta:
        author = "linted"
        info = "Weird password dump script thing"

    strings:
        $lastpast = /lastpast https:\/\/pastebin.com/
        $pastlog = /PAST LOG - https:\/\/pastebin.com/
    condition:
        $lastpast or $pastlog

}
