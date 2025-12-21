module Language.While.Lexer
    ( identifier
    , reserved
    , reservedOp
    , parens
    , integer
    , semi
    , whiteSpace
    ) where

import Text.Parsec
import Text.Parsec.String (Parser)
import Text.Parsec.Language (emptyDef)
import qualified Text.Parsec.Token as Token

languageDef :: Token.LanguageDef ()
languageDef = emptyDef
    { Token.commentLine     = "//"
    , Token.identStart      = letter
    , Token.identLetter     = alphaNum <|> char '_'
    , Token.reservedNames   = 
        [ "if", "then", "else", "end"
        , "while", "do"
        , "skip"
        , "true", "false"
        , "not", "and", "or"
        ]
    , Token.reservedOpNames = 
        [ ":="
        , "+", "-", "*", "/"
        , "=", "<", ">", "<=", ">="
        , ";"
        , "(", ")"
        ]
    }

lexer :: Token.TokenParser ()
lexer = Token.makeTokenParser languageDef

identifier :: Parser String
identifier = Token.identifier lexer

reserved :: String -> Parser ()
reserved = Token.reserved lexer

reservedOp :: String -> Parser ()
reservedOp = Token.reservedOp lexer

parens :: Parser a -> Parser a
parens = Token.parens lexer

integer :: Parser Integer
integer = Token.integer lexer

semi :: Parser String
semi = Token.semi lexer

whiteSpace :: Parser ()
whiteSpace = Token.whiteSpace lexer
