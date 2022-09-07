/*
    Financial websites
*/

rule amex
{
    meta:
        author= "@linted"
    strings:
        $americanexpress = /americanexpress/i
        $amex = /\samex/
    condition:
        any of them
}

rule chase
{
    meta:
        author = "@linted"
    strings:
        $chase = /chase.com/i
    condition:
        any of them
}

rule user_pass
{
    meta:
        author = "@linted"
    strings:
        $user = /user/i
        $pass = /pass/i
    condition:
        all of them
}