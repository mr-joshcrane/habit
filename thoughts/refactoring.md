# 27/3/2022

So I have a primitive network store using protocol buffers. Now it's time to refactor... where to start?

I think I do a lot more double handling of data that I don't need to do... 

1. Can I make the types more natural?
2. Does it make sense to store things in memory/cache?

Lets work on making the types more natural...

## bug

If you perform the same habit twice you will end up with two habits
   one correct habit 
   one habit with no name and the old timestamp