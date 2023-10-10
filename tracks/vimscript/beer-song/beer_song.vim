function! BottleString(num) abort
  if a:num > 1
    return printf("%d bottles", a:num)
  elseif a:num == 1
    return printf("%d bottle", a:num)
  else
    return "no more bottles"
  endif
endfunction

function! Verse(verse) abort
  if a:verse == 0
    return "No more bottles of beer on the wall, no more bottles of beer.\nGo to the store and buy some more, 99 bottles of beer on the wall.\n"
  else 
    return printf("%s of beer on the wall, %s of beer.\nTake %s down and pass it around, %s of beer on the wall.\n",
\       BottleString(a:verse), 
\       BottleString(a:verse), 
\       a:verse == 1 ? "it" : "one",
\       BottleString(a:verse-1))
  endif
endfunction

function! Verses(start, end) abort
  let out = Verse(a:start)
  let idx=a:start - 1
  while idx >= a:end
    let out = out . "\n" . Verse(idx)
    let idx = idx-1
  endwhile

  return out

endfunction
