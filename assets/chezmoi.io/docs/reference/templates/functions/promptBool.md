# `promptBool` *prompt* [*default*]

`promptBool` prompts the user with *prompt* and returns the user's response
interpreted as a boolean. If *default* is passed the user's response is empty
then it returns *default*. It is only available when generating the initial
config file. The user's response is interpreted as follows (case insensitive):

| Response                | Result  |
| ----------------------- | ------- |
| 1, on, t, true, y, yes  | `true`  |
| 0, off, f, false, n, no | `false` |
