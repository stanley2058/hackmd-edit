# HackMD Edit

A simple CLI tool to allow editing HackMD notes in your favorite editor via HackMD API. Currently only support single person usage, concurrent editing will cause content overwriting and data loss.

## Configuration

The tool uses the following environment variables:

- `HACKMD_API_KEY`: Your HackMD API key.
- `EDITOR`: The editor to use.
- `EDITOR_ARGS`: Arguments to pass to the editor.
  - HackMD Edit will call your editor like this: `$EDITOR $EDITOR_ARGS $FILE_TO_EDIT`

You can also use `.env` file to set these variables, see [.env.example](.env.example) for an example.

## Usage

To edit a personal note:

```
> hackmd-edit --note <note-id>
```

To edit a team note:

```
> hackmd-edit --team <team-path> --note <note-id>
```

To edit with VSCode:

```
> EDITOR=code EDITOR_ARGS=--wait hackmd-edit --note <note-id>
```
