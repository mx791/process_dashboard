# Process Dashboard

Create a visual dashboard for your long Golang processes.

![](./figures/dashboard.PNG)

## How to use
```
d := process_dashboard.DashBoard{yourTask, itersCount, []process_dashboard.Callback{
    process_dashboard.LastValueCallback{"Important metric", "metric"},
    process_dashboard.LineCallback{"Important metric evolution", "metric"},
}}
d.Run()
```


## Metrics
- LastValueCallback : Display the last raw value of a metric, as text
- LineCallback : Display the evolution of a metric in a line chart
- LineCallback : Display the evolution of a few metrics in a line chart

## Further improvement
- Allow to choose dashboard's port
- Allow to choose page title
- Provide information about task state (running, done)
- Multiline chart ✅
- Scatter plot
- Task percent complete
- Display raw text in dashboard
