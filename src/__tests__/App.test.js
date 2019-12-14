import { render } from '@testing-library/svelte';
import App from '../App.svelte';
import '@testing-library/jest-dom/extend-expect';

test('shows proper heading when rendered', () => {
  const { getByText } = render(App);

  expect(getByText('cat')).toBeInTheDocument();
  expect(getByText('dog')).toBeInTheDocument();
  expect(getByText('kek')).toBeInTheDocument();
});
