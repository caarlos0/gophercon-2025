Set FontSize 20
Set Width 1200
Set Height 600
Set BorderRadius 40
Set WindowBar Colorful
Set WindowBarSize 50
Output wish.gif

Hide
Type "tmux kill-session -t demo" Enter
Type "tmux new -A -s demo" Enter
Type "tmux set status off" Enter
Type "bash --rcfile rc.bash" Enter
Type "tmux splitw bash --rcfile rc.bash" Enter
Ctrl+B
Type k
Type "clear" Enter
Show

Type "go run ./wish/main.go"
Enter
Sleep 1s

Ctrl+B
Type j

Sleep 1s
Type "ssh -p 23234 localhost" Enter
Sleep 10s
Enter
Sleep 10s


Hide
Type "tmux kill-session -t demo" Enter
Show
