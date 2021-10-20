if exists("b:current_syntax")
  finish
endif


syn case match

syn keyword     magpieDirective         import
syn keyword     magpieDeclaration       let const


hi def link     magpieDirective         Statement
hi def link     magpieDeclaration       Type

" Linq Keywords
syn keyword     magpieLinq              from where select group into orderby join in on equals by ascending descending

syn keyword     magpieStatement         return let const spawn defer struct enum using async await service
syn keyword     magpieException         try catch finally throw
syn keyword     magpieConditional       if else elif unless and or case is
syn keyword     magpieRepeat            do while for break continue grep map
syn keyword     magpieBranch            break continue
syn keyword     magpieClass             class new property get set default this parent static public private protected interface

hi def link     magpieStatement         Statement
hi def link     magpieClass             Statement
hi def link     magpieConditional       Conditional
hi def link     magpieBranch            Conditional
hi def link     magpieLabel             Label
hi def link     magpieRepeat            Repeat
hi def link     magpieLinq              Keyword

syn match       magpieDeclaration       /\<fn\>/
syn match       magpieDeclaration       /^fn\>/


syn keyword magpieCommentTodo contained TODO FIXME XXX NOTE
hi def link magpieCommentTodo Todo

syn match comment "#.*$"    contains=@Spell,magpieCommentTodo
syn match comment "\/\/.*$" contains=@Spell,magpieCommentTodo

syn keyword     magpieCast              int str float array


hi def link     magpieCast              Type


syn keyword     magpieBuiltins          len
syn keyword     magpieBuiltins          println print stdin stdout stderr
syn keyword     magpieBoolean           true false
syn keyword     magpieNull              nil

hi def link     magpieBuiltins          Keyword
hi def link     magpieNull              Keyword
hi def link     magpieBoolean           Boolean


" Comments; their contents
syn keyword     magpieTodo              contained TODO FIXME XXX BUG
syn cluster     magpieCommentGroup      contains=magpieTodo
syn region      magpieComment           start="#" end="$" contains=@magpieCommentGroup,@Spell,@magpieTodo


hi def link     magpieComment           Comment
hi def link     magpieTodo              Todo


" magpie escapes
syn match       magpieEscapeOctal       display contained "\\[0-7]\{3}"
syn match       magpieEscapeC           display contained +\\[abfnrtv\\'"]+
syn match       magpieEscapeX           display contained "\\x\x\{2}"
syn match       magpieEscapeU           display contained "\\u\x\{4}"
syn match       magpieEscapeBigU        display contained "\\U\x\{8}"
syn match       magpieEscapeError       display contained +\\[^0-7xuUabfnrtv\\'"]+


hi def link     magpieEscapeOctal       magpieSpecialString
hi def link     magpieEscapeC           magpieSpecialString
hi def link     magpieEscapeX           magpieSpecialString
hi def link     magpieEscapeU           magpieSpecialString
hi def link     magpieEscapeBigU        magpieSpecialString
hi def link     magpieSpecialString     Special
hi def link     magpieEscapeError       Error
hi def link     magpieException		Exception

" Strings and their contents
syn cluster     magpieStringGroup       contains=magpieEscapeOctal,magpieEscapeC,magpieEscapeX,magpieEscapeU,magpieEscapeBigU,magpieEscapeError
syn region      magpieString            start=+"+ skip=+\\\\\|\\"+ end=+"+ contains=@magpieStringGroup
syn region      magpieRegExString       start=+/[^/*]+me=e-1 skip=+\\\\\|\\/+ end=+/\s*$+ end=+/\s*[;.,)\]}]+me=e-1 oneline
syn region      magpieRawString         start=+`+ end=+`+

hi def link     magpieString            String
hi def link     magpieRawString         String
hi def link     magpieRegExString       String

" Characters; their contents
syn cluster     magpieCharacterGroup    contains=magpieEscapeOctal,magpieEscapeC,magpieEscapeX,magpieEscapeU,magpieEscapeBigU
syn region      magpieCharacter         start=+'+ skip=+\\\\\|\\'+ end=+'+ contains=@magpieCharacterGroup


hi def link     magpieCharacter         Character


" Regions
syn region      magpieBlock             start="{" end="}" transparent fold
syn region      magpieParen             start='(' end=')' transparent


" Integers
syn match       magpieDecimalInt        "\<\d\+\([Ee]\d\+\)\?\>"
syn match       magpieHexadecimalInt    "\<0x\x\+\>"
syn match       magpieOctalInt          "\<0\o\+\>"
syn match       magpieOctalError        "\<0\o*[89]\d*\>"


hi def link     magpieDecimalInt        Integer
hi def link     magpieHexadecimalInt    Integer
hi def link     magpieOctalInt          Integer
hi def link     Integer                 Number

" Floating point
syn match       magpieFloat             "\<\d\+\.\d*\([Ee][-+]\d\+\)\?\>"
syn match       magpieFloat             "\<\.\d\+\([Ee][-+]\d\+\)\?\>"
syn match       magpieFloat             "\<\d\+[Ee][-+]\d\+\>"


hi def link     magpieFloat             Float
"hi def link     magpieImaginary         Number


if exists("magpie_fold")
    syn match	magpieFunction	"\<fn\>"
    syn region	magpieFunctionFold	start="\<fn\>.*[^};]$" end="^\z1}.*$" transparent fold keepend

    syn sync match magpieSync	grouphere magpieFunctionFold "\<fn\>"
    syn sync match magpieSync	grouphere NONE "^}"

    setlocal foldmethod=syntax
    setlocal foldtext=getline(v:foldstart)
else
    syn keyword magpieFunction	fn
    syn match	magpieBraces	"[{}\[\]]"
    syn match	magpieParens	"[()]"
endif

syn sync fromstart
syn sync maxlines=100

hi def link magpieFunction		Function
hi def link magpieBraces		Function

let b:current_syntax = "magpie"
