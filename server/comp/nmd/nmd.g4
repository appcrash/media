/* node manipulation & description language */
grammar nmd;

graph : stmt_list ;
stmt_list : stmt (';' stmt)* ';'? ;
stmt : node_def
     | link_stmt
     | call_stmt
     | cast_stmt
     | sink_stmt
     ;

link_stmt : endpoint '->' endpoint ('->' endpoint)* ;
call_stmt : node_def '<->' cmd=QUOTED_STRING ;
cast_stmt : node_def '<--' cmd=QUOTED_STRING ;
sink_stmt : '<-chan' channel=ID ;

endpoint : node_def
         | '{' node_def (',' node_def)* '}'
         ;

node_def : '[' node_id node_prop* ']' ;
node_id : name=ID  ('@' scope=ID)? (':' typ=ID)? ;
node_prop : key=ID '=' value=property ;


property : QUOTED_STRING  #PropQuoteString
         | ID             #PropId
         | INT            #PropInt
         | FLOAT          #PropFloat
         ;

QUOTED_STRING : '\'' ( ESC | . )*? '\'' ;
fragment
ESC : '\\\''  ;

INT : DIGIT+
    | '0x' DIGIT+
    | '0X' DIGIT+
    ;
FLOAT : DIGIT+ '.' DIGIT*
      | '.' DIGIT+
      ;

WS : [ \t\r\n]+ -> skip ;
fragment
DIGIT : [0-9] ;
fragment
LETTER : [a-zA-Z_] ;

ID : LETTER (LETTER | DIGIT)* ;

