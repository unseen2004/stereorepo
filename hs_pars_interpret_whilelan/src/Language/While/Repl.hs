module Language.While.Repl
    ( repl
    , runFile
    , runString
    ) where

import System.IO (hFlush, stdout)
import qualified Data.Map.Strict as Map
import Control.Monad (unless)

import Language.While.Parser
import Language.While.Interpreter

repl :: IO ()
repl = do
    putStrLn "While Language Interpreter v0.1"
    putStrLn "Type ':help' for commands, ':quit' to exit"
    putStrLn ""
    replLoop emptyState

replLoop :: State -> IO ()
replLoop state = do
    putStr "while> "
    hFlush stdout
    input <- getLine
    unless (null input) $ processInput state input

processInput :: State -> String -> IO ()
processInput state input
    | input == ":quit" || input == ":q" = 
        putStrLn "Goodbye!"
    
    | input == ":help" || input == ":h" = do
        printHelp
        replLoop state
    
    | input == ":state" || input == ":s" = do
        printState state
        replLoop state
    
    | input == ":clear" || input == ":c" = do
        putStrLn "State cleared."
        replLoop emptyState
    
    | take 6 input == ":load " = do
        let filename = drop 6 input
        result <- runFile filename
        case result of
            Left err -> do
                putStrLn $ "Error: " ++ show err
                replLoop state
            Right newState -> do
                putStrLn "File loaded successfully."
                printState newState
                replLoop newState
    
    | otherwise = do
        case parseString input of
            Left err -> do
                putStrLn $ "Parse error: " ++ show err
                replLoop state
            Right stmt -> do
                let newState = exec stmt state
                printState newState
                replLoop newState

printHelp :: IO ()
printHelp = do
    putStrLn "Commands:"
    putStrLn "  :help, :h     - Show this help"
    putStrLn "  :quit, :q     - Exit the REPL"
    putStrLn "  :state, :s    - Show current state"
    putStrLn "  :clear, :c    - Clear state"
    putStrLn "  :load <file>  - Load and run a file"
    putStrLn ""
    putStrLn "Example programs:"
    putStrLn "  x := 5"
    putStrLn "  x := x + 1"
    putStrLn "  if x > 0 then y := 1 else y := 0 end"
    putStrLn "  while x > 0 do x := x - 1 end"

printState :: State -> IO ()
printState state
    | Map.null state = putStrLn "  (empty state)"
    | otherwise = mapM_ printVar (Map.toList state)
  where
    printVar (name, value) = putStrLn $ "  " ++ name ++ " = " ++ show value

runFile :: FilePath -> IO (Either String State)
runFile path = do
    result <- parseFile path
    case result of
        Left err   -> return $ Left (show err)
        Right stmt -> return $ Right (execProgram stmt)

runString :: String -> Either String State
runString input = case parseString input of
    Left err   -> Left (show err)
    Right stmt -> Right (execProgram stmt)
