module Language.While.Interpreter
    ( State
    , emptyState
    , evalA
    , evalB
    , exec
    , execProgram
    , lookupVar
    ) where

import qualified Data.Map.Strict as Map
import Data.Map.Strict (Map)

import Language.While.Syntax

type State = Map String Integer

emptyState :: State
emptyState = Map.empty

lookupVar :: String -> State -> Integer
lookupVar name state = Map.findWithDefault 0 name state

updateVar :: String -> Integer -> State -> State
updateVar = Map.insert

evalA :: AExpr -> State -> Integer
evalA expr state = case expr of
    IntLit n    -> n
    Var name    -> lookupVar name state
    Add e1 e2   -> evalA e1 state + evalA e2 state
    Sub e1 e2   -> evalA e1 state - evalA e2 state
    Mul e1 e2   -> evalA e1 state * evalA e2 state
    Div e1 e2   -> evalA e1 state `div` evalA e2 state

evalB :: BExpr -> State -> Bool
evalB expr state = case expr of
    BoolLit b   -> b
    Not e       -> not (evalB e state)
    And e1 e2   -> evalB e1 state && evalB e2 state
    Or e1 e2    -> evalB e1 state || evalB e2 state
    Eq e1 e2    -> evalA e1 state == evalA e2 state
    Lt e1 e2    -> evalA e1 state <  evalA e2 state
    Gt e1 e2    -> evalA e1 state >  evalA e2 state
    Le e1 e2    -> evalA e1 state <= evalA e2 state
    Ge e1 e2    -> evalA e1 state >= evalA e2 state

exec :: Stmt -> State -> State
exec stmt state = case stmt of
    Assign name expr -> 
        updateVar name (evalA expr state) state
    
    Skip -> state
    
    Seq s1 s2 -> 
        let state' = exec s1 state
        in exec s2 state'
    
    If cond thenBranch elseBranch ->
        if evalB cond state
            then exec thenBranch state
            else exec elseBranch state
    
    While cond body ->
        if evalB cond state
            then let state' = exec body state
                 in exec (While cond body) state'
            else state

execProgram :: Stmt -> State
execProgram stmt = exec stmt emptyState
