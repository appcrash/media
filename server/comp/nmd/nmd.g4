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

link_stmt : endpoint link_operator endpoint (link_operator endpoint)* ;
call_stmt : node_def '<->' cmd=QUOTED_STRING ;
cast_stmt : node_def '<--' cmd=QUOTED_STRING ;
sink_stmt : '<-chan' node=ID ; /* pull data from node in graph */

endpoint : node_def
         | '{' node_def (',' node_def)* '}'
         ;

node_def : '[' node_id node_prop* ']'  ;
node_id : name=ID  ('@' scope=ID)? (':' typ=ID)? ;
node_prop : key=ID '=' value=property ;

msg_type_list :  ID (',' ID)*  ;
link_operator : '<' msg_type_list '>'
              | '->' ;

property : QUOTED_STRING  #PropQuoteString
         | ID             #PropId
         | INT            #PropInt
         | FLOAT          #PropFloat
         ;

QUOTED_STRING : '\'' ( ESC | . )*? '\'' ;
fragment
ESC : '\\\''  ;

INT : DIGIT+
    | '0x' HEXDIGIT+
    | '0X' HEXDIGIT+
    ;
FLOAT : DIGIT+ '.' DIGIT*
      | '.' DIGIT+
      ;

WS : [ \t\r\n]+ -> skip ;
fragment
DIGIT : [0-9] ;
fragment
HEXDIGIT : [0-9a-f] ;
fragment
LETTER : [a-zA-Z_] ;

ID : LETTER (LETTER | DIGIT)* ;

