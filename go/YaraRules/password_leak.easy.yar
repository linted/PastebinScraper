/*
    These rules attempt to find password leaks / dumps
*/

rule email_list
{
    meta:
        author = "@KevTheHermit"
        info = "Part of PasteHunter"
        reference = "https://github.com/kevthehermit/PasteHunter"

    strings:
        $email_add = /\b[\w\.-]+@[\w\.-]+\.\w+\b/
        $email = /\b((([!#$%&'*+\-\/=?^`{|}~\w])|([!#$%&'*+\-\/=?^`{|}~\w][!#$%&'*+\-\/=?^`{|}~\.\w]{0,}[!#$%&'*+\-\/=?^`{|}~\w]))[@]\w+([-.]\w+)*\.\w+([-.]\w+)*)\b/
        $extm3u = /#EXTM3U/
        $etxinf = /#ETXINF/
        $extinf = /#EXTINF/
    condition:
        (#email_add >= 5 or #email >= 4) and not ($extm3u or $etxinf or $extinf)

}


/*
rule password_list
{
    meta:
        author = "@KevTheHermit"
        info = "Part of PasteHunter"
        reference = "https://github.com/kevthehermit/PasteHunter"

    strings:
        $data_format = /\b([@a-zA-Z0-9._-]{5,})(:|\|)(.*)\b/

    condition:
        #data_format >= 10

}
*/