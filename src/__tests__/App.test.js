import { render } from '@testing-library/svelte';
import App from '../App.svelte';
import '@testing-library/jest-dom/extend-expect';

test('shows proper heading when rendered', () => {
  const { getByText } = render(App);

  expect(getByText('Hello new world!')).toBeInTheDocument();
});
