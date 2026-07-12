import type { PageLoad } from './$types';
import { GetWorkbook } from '../../../bindings/merger/services/workbook';
import type { WorkbookInfo } from '../../../bindings/merger/utility/models';

export const load: PageLoad = async ({ url }) => {
	return await GetWorkbook(url.searchParams.get('id'));
};
