import './css/DataGrid.css';
import classNames from 'classnames/bind';
import PropTypes from 'prop-types';
import Column from '../common/Column.jsx';
import Columns from '../common/Columns.jsx';

/**
 * 数据表格
 * @author tengge / https://github.com/tengge1
 */
class DataGrid extends React.Component {
    constructor(props) {
        super(props);

        this.handleClick = this.handleClick.bind(this, props.onSelect);
    }

    render() {
        const { className, style, children, data, pageSize, pageNum, total, selected } = this.props;

        const columns = children.props.children.map(n => {
            return {
                type: n.props.type,
                field: n.props.field,
                title: n.props.title,
            };
        });

        const header = <thead>
            <tr>
                {columns.map(n => {
                    let field = n.type === 'number' ? 'number' : n.field;
                    return <td name={n.field} key={field}>{n.title}</td>;
                })}
            </tr>
        </thead>;

        const body = <tbody>
            {data.map((n, i) => {
                return <tr className={selected === n.id ? 'selected' : null} data-id={n.id} key={n.id} onClick={this.handleClick}>
                    {columns.map((m, j) => {
                        if (m.type === 'number') {
                            return <td key={'number'}>{i + 1}</td>;
                        } else {
                            return <td key={m.field}>{n[m.field]}</td>;
                        }
                    })}
                </tr>;
            })}
        </tbody>;

        return <div className={classNames('DataGrid', className)} style={style}>
            <table>
                {header}
                {body}
            </table>
        </div>;
    }

    handleClick(onSelect, event) {
        const id = event.currentTarget.getAttribute('data-id');

        const record = this.props.data.filter(n => n.id === id)[0];

        onSelect && onSelect(record);
    }
}

DataGrid.propTypes = {
    className: PropTypes.string,
    style: PropTypes.object,
    children: (props, propName, componentName) => {
        const children = props[propName];
        if (children.type !== Columns) {
            return new TypeError(`Invalid prop \`${propName}\` of type \`${children.type.name}\` supplied to \`${componentName}\`, expected \`Columns\`.`);
        }
    },
    pages: PropTypes.bool,
    data: PropTypes.array,
    pageSize: PropTypes.number,
    pageNum: PropTypes.number,
    total: PropTypes.number,
    selected: PropTypes.string,
    onSelect: PropTypes.func,
};

DataGrid.defaultProps = {
    className: null,
    style: null,
    children: null,
    pages: false,
    data: [],
    pageSize: 20,
    pageNum: 1,
    total: 0,
    selected: null,
    onSelect: null,
};

export default DataGrid;