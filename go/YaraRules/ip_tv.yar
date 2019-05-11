/*
    These rules identify ip tv listings
*/

rule ip_tv
{
    meta:
        author = "@linted"
    strings:
        $extm3u = /#EXTM3U/
        $etxinf = /#ETXINF/
        $extinf = /#EXTINF/
    condition:
        any of them
}
