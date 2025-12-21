module Language.While.Parser
    ( parseProgram
    , parseString
    , parseFile
    ) where

import Text.Parsec
import Text.Parsec.String (Parser, parseFromFile)
import Text.Parsec.Expr
import Data.Functor.Identity (Identity)

import Language.While.Syntax
import Language.While.Lexer

program :: Parser Stmt
program = do
    whiteSpace
    stmt <- statement
    eof
    return stmt

statement :: Parser Stmt
statement = chainr1 singleStatement seqOp
  where
    seqOp = do
        reservedOp ";"
        return Seq

singleStatement :: Parser Stmt
singleStatement =  ifStmt
               <|> whileStmt
               <|> skipStmt
               <|> assignStmt

skipStmt :: Parser Stmt
skipStmt = do
    reserved "skip"
    return Skip

assignStmt :: Parser Stmt
assignStmt = do
    var <- identifier
    reservedOp ":="
    expr <- aExpr
    return $ Assign var expr

ifStmt :: Parser Stmt
ifStmt = do
    reserved "if"
    cond <- bExpr
    reserved "then"
    thenBranch <- statement
    reserved "else"
    elseBranch <- statement
    reserved "end"
    return $ If cond thenBranch elseBranch

whileStmt :: Parser Stmt
whileStmt = do
    reserved "while"
    cond <- bExpr
    reserved "do"
    body <- statement
    reserved "end"
    return $ While cond body

aExpr :: Parser AExpr
aExpr = buildExpressionParser aOperators aTerm

aOperators :: [[Operator String () Identity AExpr]]
aOperators =
    [ [Infix (reservedOp "*" >> return Mul) AssocLeft,
       Infix (reservedOp "/" >> return Div) AssocLeft]
    , [Infix (reservedOp "+" >> return Add) AssocLeft,
       Infix (reservedOp "-" >> return Sub) AssocLeft]
    ]

aTerm :: Parser AExpr
aTerm =  parens aExpr
     <|> (IntLit <$> integer)
     <|> (Var <$> identifier)

bExpr :: Parser BExpr
bExpr = buildExpressionParser bOperators bTerm

bOperators :: [[Operator String () Identity BExpr]]
bOperators =
    [ [Prefix (reserved "not" >> return Not)]
    , [Infix (reserved "and" >> return And) AssocLeft]
    , [Infix (reserved "or"  >> return Or)  AssocLeft]
    ]

bTerm :: Parser BExpr
bTerm =  parens bExpr
     <|> (reserved "true"  >> return (BoolLit True))
     <|> (reserved "false" >> return (BoolLit False))
     <|> relExpr

relExpr :: Parser BExpr
relExpr = do
    a1 <- aExpr
    op <- relOp
    a2 <- aExpr
    return $ op a1 a2

relOp :: Parser (AExpr -> AExpr -> BExpr)
relOp =  (reservedOp "="  >> return Eq)
     <|> (reservedOp "<=" >> return Le)
     <|> (reservedOp ">=" >> return Ge)
     <|> (reservedOp "<"  >> return Lt)
     <|> (reservedOp ">"  >> return Gt)

parseProgram :: Parser Stmt
parseProgram = program

parseString :: String -> Either ParseError Stmt
parseString = parse program ""

parseFile :: FilePath -> IO (Either ParseError Stmt)
parseFile = parseFromFile program
