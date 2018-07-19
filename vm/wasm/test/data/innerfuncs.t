
// think how to import not from "env"   namespaces ???

extern "C" void print(int v);

int test() {
  print (10);
  return 20;
}

===

(module
 (type $FUNCSIG$vi (func (param i32)))
 (import "INS" "print2" (func $print (param i32)))
 (table 0 anyfunc)
 (memory $0 1)
 (export "memory" (memory $0))
 (export "_Z4testv" (func $_Z4testv))
 (func $_Z4testv (; 1 ;) (result i32)
  (call $print
   (i32.const 10)
  )
  (i32.const 20)
 )
)

===

[
{
 "name" : "Just call",
 "fidx": 1,
 "args": [],
 "global": [],
 "gas": 14,
 "result": [0,0,0,10],
 "memcheck": {"12":10}
}
]
