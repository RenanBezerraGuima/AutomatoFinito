# AutomatoFinito Geral

## Overview

This program simulates Finite Automata, which are fundamental concepts in computer science and formal language theory. It supports both:

*   **Deterministic Finite Automata (DFAs)**
*   **Non-deterministic Finite Automata (NFAs)**, including those with **epsilon (ε) transitions**.

The primary purpose of this tool is for educational use, allowing users to design, visualize the structure of, and test various automata.

## Features

*   Creation and simulation of DFAs.
*   Creation and simulation of NFAs, including support for ε-transitions.
*   Interactive console-based user interface for defining automata components (states, alphabet, transitions, etc.).
*   Testing of input strings against the currently defined automaton.
*   Pre-defined examples to demonstrate DFA and NFA functionalities.
*   Input validation to guide the user and prevent common errors during automaton definition.

## Getting Started

### Prerequisites

*   Go programming language environment (Go 1.18 or later is recommended). You can download it from [https://golang.org/dl/](https://golang.org/dl/).

### Compilation

To compile the program, navigate to the directory containing the source file (`automatoFinitoGeral.go`) and run:

```bash
go build automatoFinitoGeral.go
```

### Running the Program

After successful compilation, you can run the program using:

*   On Linux/macOS:
    ```bash
    ./automatoFinitoGeral
    ```
*   On Windows:
    ```bash
    automatoFinitoGeral.exe
    ```

## How to Use

Upon running the program, you will be greeted with a main menu:

```
Menu Principal:
1. Rodar Exemplo Pré-definido
2. Criar Novo Autômato
3. Sair
Escolha uma opção:
```

*   **1. Rodar Exemplo Pré-definido:** Shows the execution of built-in examples, including an NFA that accepts strings ending with "ab" and an NFA using epsilon transitions for the language "a*b".
*   **2. Criar Novo Autômato:** Allows you to define your own automaton step-by-step.
*   **3. Sair:** Exits the program.

### Defining an Automaton

If you choose to create a new automaton, the program will guide you through several steps:

1.  **States:**
    *   Prompt: `Digite os estados (um por vez). Digite "fim" para encerrar:`
    *   Input each state name on a new line.
    *   Type `fim` to finish adding states.
    *   Constraints: State names cannot be empty or duplicated.

2.  **Alphabet:**
    *   Prompt: `Digite o alfabeto (um símbolo por vez). Digite "fim" para encerrar:`
    *   Input each alphabet symbol on a new line.
    *   Type `fim` to finish adding symbols.
    *   Constraints: Symbols must be single characters. Duplicated symbols are not allowed.
    *   Note: The epsilon symbol (`ε`) is handled implicitly for transitions; do not add it to the alphabet here unless it's a regular symbol in your language (which is rare for standard automata definitions).

3.  **Initial State:**
    *   Prompt: `Digite o estado inicial:`
    *   Enter the name of one of the states you defined earlier.
    *   Constraint: The state must exist in your list of defined states.

4.  **Transitions:**
    This is where you define how the automaton moves from one state to another based on input symbols.
    *   The process starts with: `--- Adicionar Transições ---`
    *   **Origin State:**
        *   Prompt: `Origem (ou "fim" para encerrar tudo):`
        *   Enter the name of the state from which the transition originates. Type `fim` to finish defining all transitions.
    *   **Symbol:**
        *   Prompt: `Símbolo para {origem} (ou "eps"/"epsilon" para épsilon):`
        *   Enter the input symbol for the transition.
        *   For an **epsilon transition**, type `eps` or `epsilon`.
    *   **Destination State(s):**
        *   Prompt: `Adicionando destinos para ({origem}, '{simbolo}'):`
        *   Prompt: `  Destino para {origem},'{simbolo}' (ou "fim" para esta transição):`
        *   Enter a destination state.
        *   For NFAs, you can enter multiple destination states for the same origin state and symbol. The program will loop, asking for a new destination state for the current (origin, symbol) pair until you type `fim`.
        *   Typing `fim` here finalizes the destinations for the *current (origin, symbol) pair* and allows you to define a new transition (new origin or new symbol).
    *   Constraints:
        *   Origin and destination states must be from your defined list of states.
        *   Symbols (if not epsilon) must be part of your defined alphabet.

5.  **Final States:**
    *   Prompt: `Digite os estados finais (um por vez). Digite "fim" para encerrar:`
    *   Input each final state name on a new line.
    *   Type `fim` to finish adding final states.
    *   Constraints: Final states must be from your defined list of states and cannot be duplicated in the list of final states.

### Testing Strings

After successfully defining an automaton:
1.  The program will display the details of the automaton you created.
2.  It will then prompt: `Digite a cadeia para testar (ou "sair" para encerrar):`
3.  Enter any string you want to test.
4.  The program will output whether the string is `aceita` (accepted) or `não aceita` (not accepted).
5.  To stop testing and return to the main menu, type `sair`.

## Example of NFA Definition

Let's define an NFA that accepts strings containing "aa" (i.e., L = {x | x contains "aa" as a substring}).
*   **States:** q0, q1, q2
*   **Alphabet:** a, b
*   **Initial State:** q0
*   **Transitions:**
    *   q0, a -> q0, q1  (Stays in q0 for other 'a's, or moves to q1 for the first 'a' of "aa")
    *   q0, b -> q0      (Stays in q0 for 'b's)
    *   q1, a -> q2      (Second 'a' of "aa" moves to final state q2)
    *   q2, a -> q2      (Optional: Stays in q2 if more 'a's follow "aa")
    *   q2, b -> q2      (Optional: Stays in q2 if 'b's follow "aa")
*   **Final States:** q2

Here's how you would input this into the program:

1.  **Main Menu:** Choose option `2` (Criar Novo Autômato).

2.  **States:**
    ```
    Digite os estados (um por vez). Digite "fim" para encerrar:
    > q0
    Estado 'q0' adicionado.
    > q1
    Estado 'q1' adicionado.
    > q2
    Estado 'q2' adicionado.
    > fim
    ```

3.  **Alphabet:**
    ```
    Digite o alfabeto (um símbolo por vez). Digite "fim" para encerrar:
    > a
    Símbolo 'a' adicionado ao alfabeto.
    > b
    Símbolo 'b' adicionado ao alfabeto.
    > fim
    ```

4.  **Initial State:**
    ```
    Digite o estado inicial: > q0
    ```

5.  **Transitions:**
    ```
    --- Adicionar Transições ---
    Para cada transição, primeiro o estado de origem e o símbolo.
    Depois, digite os estados de destino um por vez.
    Digite "fim" como estado de destino para finalizar a adição de destinos para a transição atual.
    Digite "fim" como origem para encerrar a adição de todas as transições.

    Origem (ou "fim" para encerrar tudo): > q0
    Símbolo para q0 (ou "eps"/"epsilon" para épsilon): > a
    Adicionando destinos para (q0, 'a'):
      Destino para q0,'a' (ou "fim" para esta transição): > q0
        Adicionado: q0 --'a'--> q0
      Destino para q0,'a' (ou "fim" para esta transição): > q1
        Adicionado: q0 --'a'--> q1
      Destino para q0,'a' (ou "fim" para esta transição): > fim
    Próxima transição.

    Origem (ou "fim" para encerrar tudo): > q0
    Símbolo para q0 (ou "eps"/"epsilon" para épsilon): > b
    Adicionando destinos para (q0, 'b'):
      Destino para q0,'b' (ou "fim" para esta transição): > q0
        Adicionado: q0 --'b'--> q0
      Destino para q0,'b' (ou "fim" para esta transição): > fim
    Próxima transição.

    Origem (ou "fim" para encerrar tudo): > q1
    Símbolo para q1 (ou "eps"/"epsilon" para épsilon): > a
    Adicionando destinos para (q1, 'a'):
      Destino para q1,'a' (ou "fim" para esta transição): > q2
        Adicionado: q1 --'a'--> q2
      Destino para q1,'a' (ou "fim" para esta transição): > fim
    Próxima transição.

    Origem (ou "fim" para encerrar tudo): > q2
    Símbolo para q2 (ou "eps"/"epsilon" para épsilon): > a
    Adicionando destinos para (q2, 'a'):
      Destino para q2,'a' (ou "fim" para esta transição): > q2
        Adicionado: q2 --'a'--> q2
      Destino para q2,'a' (ou "fim" para esta transição): > fim
    Próxima transição.
    
    Origem (ou "fim" para encerrar tudo): > q2
    Símbolo para q2 (ou "eps"/"epsilon" para épsilon): > b
    Adicionando destinos para (q2, 'b'):
      Destino para q2,'b' (ou "fim" para esta transição): > q2
        Adicionado: q2 --'b'--> q2
      Destino para q2,'b' (ou "fim" para esta transição): > fim
    Próxima transição.

    Origem (ou "fim" para encerrar tudo): > fim
    ```

6.  **Final States:**
    ```
    Digite os estados finais (um por vez). Digite "fim" para encerrar:
    > q2
    Estado final 'q2' adicionado.
    > fim
    ```

After this, the program will display the automaton's structure and prompt for strings to test. For example:
*   `aa` -> `aceita`
*   `baa` -> `aceita`
*   `bab` -> `não aceita`
*   `a` -> `não aceita`
*   `aabaa` -> `aceita`

## Future Enhancements (Optional)

*   NFA to DFA conversion.
*   Ability to save defined automata to a file and load them later.
*   Graphical representation of automata (e.g., using a library or exporting to DOT format).
*   Minimization of DFAs.
```
