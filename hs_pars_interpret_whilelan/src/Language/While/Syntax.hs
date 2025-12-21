module Language.While.Syntax
    ( AExpr(..)
    , BExpr(..)
    , Stmt(..)
    ) where

data AExpr
    = IntLit Integer
    | Var String
    | Add AExpr AExpr
    | Sub AExpr AExpr
    | Mul AExpr AExpr
    | Div AExpr AExpr
    deriving (Show, Eq)

data BExpr
    = BoolLit Bool
    | Not BExpr
    | And BExpr BExpr
    | Or BExpr BExpr
    | Eq AExpr AExpr
    | Lt AExpr AExpr
    | Gt AExpr AExpr
    | Le AExpr AExpr
    | Ge AExpr AExpr
    deriving (Show, Eq)

data Stmt
    = Assign String AExpr
    | Skip
    | Seq Stmt Stmt
    | If BExpr Stmt Stmt
    | While BExpr Stmt
    deriving (Show, Eq)
