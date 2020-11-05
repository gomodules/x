package strings_test

import (
	"fmt"
	"testing"

	. "gomodules.xyz/x/strings"
)

const testStr = `
Quite
things
in
a
row
and

things

not

in

a
*******************************************
********************************************

**************************************
  *************
  ******** ******  ************   **************






***************************************************
row

and

things



that




are



too





far
end`

func TestFmt(t *testing.T) {
	ans := Fmt(testStr)
	fmt.Println(ans)
}
