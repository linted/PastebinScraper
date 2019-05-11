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
    condition:
        #email_add >= 5

}

rule email_better_list
{
    meta:
        author = "@linted"
    strings:
        $email = /\b((([!#$%&'*+\-\/=?^`{|}~\w])|([!#$%&'*+\-\/=?^`{|}~\w][!#$%&'*+\-\/=?^`{|}~\.\w]{0,}[!#$%&'*+\-\/=?^`{|}~\w]))[@]\w+([-.]\w+)*\.\w+([-.]\w+)*)\b/
    condition:
        #email >= 4
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