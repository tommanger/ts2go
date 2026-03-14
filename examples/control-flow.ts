function checkNumber(n: number): string {
  if (n > 0) {
    return "positive";
  } else {
    if (n < 0) {
      return "negative";
    } else {
      return "zero";
    }
  }
}

let count: number = 0;
while (count < 5) {
  console.log(count);
  count = count + 1;
}

const result: string = checkNumber(42);
console.log(result);
