function add(a: number, b: number): number {
  return a + b;
}

function subtract(a: number, b: number): number {
  return a - b;
}

function multiply(a: number, b: number): number {
  return a * b;
}

function divide(a: number, b: number): number {
  if (b == 0) {
    return 0;
  }
  return a / b;
}

const x: number = 10;
const y: number = 5;

console.log(add(x, y));
console.log(subtract(x, y));
console.log(multiply(x, y));
console.log(divide(x, y));
