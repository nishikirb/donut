# donut

Tiny dotfiles management tool written in Go.

## Installation

Homebrew

```
brew install gleamsoda/tap/donut
```

source

```
go install github.com/gleamsoda/donut/cmd/donut@latest
```

## Usage

1. Creating a configuration file.  
   The configuration file will be created at `$HOME/.config/donut/donut.toml`.

```
donut init
```

2. Move your dotfiles to the directory managed by donut (default: `$HOME/.local/share/donut`).

```
git clone your/dotfiles.git ~/.local/share/donut
```

3. Check the list and changes of files managed by donut.

```
donut list  // displays the list of files
donut diff  // displays the changes between source and destination files
```

4. Handle the changes between source and destination files.

```
donut apply // apply the changes
donut merge // merge the changes with merge tool
```

## Configuration

The configuration file `donut.toml` can be placed in the following locations:

- `$XDG_CONFIG_HOME`
- `$XDG_CONFIG_HOME/donut`
- `$HOME/.config`
- `$HOME/.config/donut`

### Configuration Options

```toml
# 'source' is the source directory containing the files that will be managed.
source = '$HOME/.local/share/donut'
# 'destination' is the directory where the files will be applied.
destination = '$HOME'
# 'editor' is the command or executable to be used as the text editor.
editor = ["vim"]
# 'pager' is the command or executable to be used as the pager for viewing file differences.
pager = ["less", "-R"]
# 'diff' is the command or executable to be used for displaying differences between files.
diff = ["diff", "-upN", "{{.Destination}}", "{{.Source}}"]
# 'merge' is the command or executable to be used for merging file changes.
merge = ["nvim", "-d", "{{.Destination}}", "{{.Source}}"]
# 'excludes' is a list of files or directories to be excluded from management.
excludes = []
```

You can modify these configuration options according to your needs in the configuration file. Ensure that the paths and commands are correctly set to match your system.
