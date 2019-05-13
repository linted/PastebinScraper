rule zeros
{
    strings:
        $oo = $00
    condition:
        $oo at 0
}