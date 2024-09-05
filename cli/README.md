# AutoAction CLI

The CLI for AutoAction tasks: `st3llar`

## Commands

1. Add/New commands:
    - All the commands should be well-grouped by their responsibilities.
    - `cobra-cli add sub-command`, which should be executed in the `.workspace/cli/` folder, and will create a template
      file under the path of `.workspace/cli/cmd/`
    - Do remember to add the subcommand package get initialized in th entrance package: `_ "github.
   com/57blocks/auto-action/cli/internal/command/general"` in main.

## Scheduler expression types

1. Rate-based: [rate-based](https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#rate-based)  
    
2. Cron-based: [cron-based](https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#cron-based)

3. One-time: [one-time](https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#one-time)
