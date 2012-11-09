/* 
 * Sample source for course work of Database System Design, 2012.09-11, Zhejiang University.
 * Source code by Pengyu CHEN(cpy.prefers.you@gmail.com).
 * COPYLEFT, ALL WRONGS RESERVED.
 */

package main

//extern interp_init
func c_interp_init() int

//extern interp_main_loop
func c_interp_main_loop() int

func main() {
    c_interp_init();
    c_interp_main_loop();
}

