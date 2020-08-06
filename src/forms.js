export function dirty(dirties, el) {
  const makeDirty = () => {
    dirties.add(el);
  };

  const isDirty = () => dirties.has(el);

  return [makeDirty, isDirty];
}
