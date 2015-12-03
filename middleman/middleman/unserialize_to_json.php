<?php

if ($argc < 2) {
    echo '';
}

$valueObj = unserialize($argv[1]);
$encodedJSON = json_encode($valueObj);

if ($encodedJSON === false) {
    echo '';
}

echo $encodedJSON;
