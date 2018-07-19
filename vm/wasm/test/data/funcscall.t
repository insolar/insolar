
static int v;
extern "C" int f(int a) {
    (a != 0) && (v = a);
    return v;
}

===

(module
 (table 0 anyfunc)
 (memory $0 1)
 (export "memory" (memory $0))
 (export "f" (func $f))
 (func $f (; 0 ;) (param $0 i32) (result i32)
  (block $label$0
   (br_if $label$0
    (i32.eqz
     (get_local $0)
    )
   )
   (i32.store offset=12
    (i32.const 0)
    (get_local $0)
   )
   (return
    (get_local $0)
   )
  )
  (i32.load offset=12
   (i32.const 0)
  )
 )
)

===

[
{
 "name" : "Just call",
 "fidx": 0,
 "args": [10],
 "global": [],
 "gas": 14,
 "result": [0,0,0,10],
 "memcheck": {"12":10}
},
{
 "name" : "With initialized memory",
 "fidx": 0,
 "args": [0],
 "global": [],
 "gas": 14,
 "result": [0,0,0,20],
 "memory": [0,0,0,0, 0,0,0,0, 0,0,0,0, 20]
}

]
