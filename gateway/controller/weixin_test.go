package controller

import (
	"log"
	"testing"
)

func TestParseXML(t *testing.T) {
	var text = `<xml>
    <AppId><![CDATA[wx1c2695469ae47724]]></AppId>
    <Encrypt><![CDATA[kONzQNhVv3Y6xt27ngNZvmEqGvx4HpzRdd8gl9RA5QItfWsp+f1RgOZtnH7tVfsoumTP1se+gYE7yoOrghimSQq2nHVVLra4ZKqjJSilywtgMxfvOo9byc/lJPG2VZVbPTLUdtBhqMmu3KQCQ0642DYpg/9I7YKNYmj79uMb7MfPsEzeMIc2PkRiIJw1qPtsF24dSYR3Mxz2BVV3zZVKGU9YFMrym+PdlBndbixBVaqD8aZBlv+VQ2b9jLMGV5dEUAC9WCO2waXG4Y7g/RSNhBim3bIn3KDdOPmB77U4CZdHuhlUpbr206HdxaG32KVr/tz0Ja7WjZsdTCH3rq8XdEcpeSp0kJrfpYkW00tMfiLKGp1au/a0C30bPvrzMv3YuOkdHKx9gXvn296ngElHHFGbG9Vi+zhxPZhUNzOKdId///nvo5w6ZWCjP5ZKKFTfSKL7pRP3fX3zWWZagz1wjw==]]></Encrypt>
</xml>`
	log.Fatalln(text)
}
