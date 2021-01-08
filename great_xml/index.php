<?php
function getXML() {
    $in = fopen('php://stdin', 'r');
    while(!feof($in)){
        $text = $text . fgets($in, 4096);
    }
    return $text;
}



libxml_disable_entity_loader (false);
// $xmlfile = file_get_contents('php://input');
// $xmlfile = trim(fgets(STDIN));

$xmlfile = $_REQUEST['xml']; //getXML();
$dom = new DOMDocument();
$dom->loadXML($xmlfile, LIBXML_NOENT | LIBXML_DTDLOAD);
$info = simplexml_import_dom($dom);
$name = $info->name;
echo "Hello $name!\n";
?>
