module Main (main) where

import System.Environment (getArgs)
import System.Exit (exitFailure)
import qualified Data.Map.Strict as Map

import Language.While.Repl (repl, runFile)

main :: IO ()
main = do
    args <- getArgs
    case args of
        []     -> repl
        [file] -> runFileAndPrint file
        _      -> printUsage

runFileAndPrint :: FilePath -> IO ()
runFileAndPrint path = do
    result <- runFile path
    case result of
        Left err -> do
            putStrLn $ "Error: " ++ err
            exitFailure
        Right state -> do
            putStrLn "Program executed successfully."
            putStrLn "Final state:"
            if Map.null state
                then putStrLn "  (no variables)"
                else mapM_ printVar (Map.toList state)
  where
    printVar (name, value) = putStrLn $ "  " ++ name ++ " = " ++ show value

printUsage :: IO ()
printUsage = do
    putStrLn "While Language Interpreter v0.1"
    putStrLn ""
    putStrLn "Usage:"
    putStrLn "  while-lang              Start interactive REPL"
    putStrLn "  while-lang <file.while> Run a file"
