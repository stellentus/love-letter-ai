# Love Letter AI

The goal of this project is to create a simple RL AI for 2-player Love Letter ([Love Letter Rules PDF](http://alderac.com/wp-content/uploads/2017/11/Love-Letter-Premium_Rulebook.pdf)). The project exists to practice implementing basic RL agents in go. Future work will target variations of the existing agents, the ability to save and load trained agents, and possibly a way to play against the agents.

There is a Monte Carlo agent in the `montecarlo` package and Sarsa in the `td` package. The code can run with the following commands (all in the `cmd` directory):
* `randomfight`: Play two (biased) random players against each other. This shows what win rate to expect for the starting player (0) compared to the second player (1). (The expected rate is about 1000-915, or 52%, showing a small advantage from starting.)
* `mcfight`: First, train MC against the (biased) random player. Then, train MC against itself in 5 rounds with epsilon decreasing each time. Finally, play greedily against random to test performance.
* `sarsafight`: Train Sarsa against the (biased) random player, with decreasing alpha and epsilon. Then play against random to test performance.

The `rules` package contains structures for the deck, allowed actions, and the game state. The `gamemaster` package can be used to run a series of games. It can also provide a trace of actions that were taken in a game. The `state` package converts game states, actions, and state-action pairs into integers for indexing. Some game state is compressed (i.e. the complete history of card plays and each player's potential knowledge of opponents' cards).

The `players` package contains a structure for simplified state (to reduce complexity), similar to the code in `state`. It also contains a biased random player. This player will randomly choose a `rules.Action`. However, if the choice is guaranteed to result in a loss, it will not be chosen (if a non-loss choice is available). This is to avoid wasting training time on obviously bad choices.
