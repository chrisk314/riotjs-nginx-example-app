const path = require('path')
const CompressionPlugin = require('compression-webpack-plugin')
const CopyWebpackPlugin = require('copy-webpack-plugin')

module.exports = {
  entry: './src/main.js',
  output: {
    path: path.resolve(__dirname, 'dist'),
    publicPath: '/dist/',
    filename: 'bundle.js'
  },
  plugins: [
    new CopyWebpackPlugin([
      { from: './src/index.html', to: './index.html' },
      { from: './src/assets', to: './assets' },
    ]),
    new CompressionPlugin({
      filename: '[path].gz[query]',
      algorithm: 'gzip',
      test: /\.(js|css|html|svg)$/,
      compressionOptions: { level: 9 },
      threshold: 256,
      minRatio: 0.8,
      deleteOriginalAssets: false,
    })
  ],
  module: {
    rules: [
      {
        test: /\.riot$/,
        exclude: /node_modules/,
        use: [{
          loader: '@riotjs/webpack-loader',
          options: {
            hot: true
          }
        }]
      },
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: ['@babel/preset-env']
          }
        }
      }
    ]
  },
  devServer: {
    clientLogLevel: 'debug',
    contentBase: path.join(__dirname, 'src'),
    compress: true,
    disableHostCheck: true,
    hot: true,
    host: "0.0.0.0",
    inline: true,
    port: "3000",
    watchContentBase: true
  }
}
