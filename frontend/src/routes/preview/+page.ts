import type { PageLoad } from './$types';
import { GetWorkbook } from '../../../bindings/merger/services/workbook';

export const load: PageLoad = async ({ url }) => {
	return await GetWorkbook(url.searchParams.get('id') ?? "");
};
