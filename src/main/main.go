/* 
 * Sample source for course work of Database System Design, 2012.09-11, Zhejiang University.
 * Source code by Pengyu CHEN(cpy.prefers.you@gmail.com).
 * COPYLEFT, ALL WRONGS RESERVED.
 */

package main

//import "../bufman"
import "../catman"
import "../common"
import "../core"
import "../idxman"
import "../interp"
import "../recman"

func c_interp_init() int __asm__("interp_init");

func c_interp_main_loop() int __asm__("interp_main_loop");

// since gccgo use a recursive init from main..
func fake_func_init() {
    catman.AllTables() 
    common.FakeInit()
    core.DropTable("")
    idxman.NewEmpty(0, 0)
    go_c_wrapper.DropTableCallback(nil)
    recman.DeleteAll(nil, nil)
    return
}

func main() {
    c_interp_init()
    c_interp_main_loop()
    return
    fake_func_init()
}

