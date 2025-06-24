const BeanShape: React.FC<{
  position: string;
  color: string;
  style?: React.CSSProperties;
}> = ({ position, color, style }) => {
  return (
    <div className={`absolute ${position}`} style={style}>
      <svg
        className="w-64 h-64" // Give explicit width/height for debugging
        viewBox="0 0 200 200"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          d="M43.2,-57.1C55.7,-45.2,65.4,-31.1,71.7,-14.5C78,2.2,80.8,21.3,73.5,34.3C66.2,47.2,48.8,54,31.2,61.5C13.6,68.9,-4.2,77,-19.9,74.1C-35.7,71.2,-49.4,57.2,-60.1,41.2C-70.8,25.1,-78.5,6.9,-76.6,-10.3C-74.7,-27.5,-63.2,-43.6,-48.8,-55.9C-34.3,-68.2,-17.2,-76.9,-0.2,-76.6C16.7,-76.4,33.4,-67.3,43.2,-57.1Z"
          transform="translate(100 100)"
          fill={color}
        />
      </svg>
    </div>
  );
};

export default BeanShape;
