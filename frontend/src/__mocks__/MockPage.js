export function mockPageFunction() {
    return true;
}

export function navigate(url) {
    return `Navigated to ${url}`;
}

export function checkInput(input) {
    return input && input.trim() !== '';
}