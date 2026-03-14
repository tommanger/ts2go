function fibonacci(n: number): number {
  if (n <= 1) {
    return n;
  }
  return fibonacci(n - 1) + fibonacci(n - 2);
}

let i: number = 0;
for (i = 0; i < 10; i = i + 1) {
  console.log(fibonacci(i));
}
