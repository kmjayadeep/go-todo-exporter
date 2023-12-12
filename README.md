# go-todo-exporter

A prometheus exporter to expose metrics from my todo list in taskwarrior

available metrics:

```
Gauge: taskwarrior_todos
```

labels:
```
status: done
status: pending
```

# Nix flake

This package can be installed in nixos using flakes
