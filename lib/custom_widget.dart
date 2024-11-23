import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

class CustomButton extends StatelessWidget {
  const CustomButton({
    super.key,
    required this.text,
    required this.onPressed,
    this.backgroundColor = Colors.white,
    this.size = const Size(double.infinity, 70),
    this.padding = const EdgeInsets.symmetric(),
    this.textPaddingLeft = 0,
    this.image,
  });

  final Color backgroundColor;
  final AutoSizeText text;
  final VoidCallback? onPressed;
  final Size size;
  final EdgeInsets padding;
  final double textPaddingLeft;
  final Image? image;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: padding,
      child: ElevatedButton(
        style: ElevatedButton.styleFrom(
          backgroundColor: backgroundColor,
          elevation: 3,
          minimumSize: size,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(32.0),
          ),
        ),
        onPressed: onPressed,
        child: Stack(
          alignment: Alignment.center,
          children: [
            if (image != null)
              Positioned(
                left: 0,
                child: image!,
              ),
            Padding(padding: EdgeInsets.only(left: textPaddingLeft), child: text),
            Center(),
          ],
        ),
      ),
    );
  }
}


class CustomClassicForm extends StatelessWidget {
  const CustomClassicForm({super.key, required this.hint,});

  final String hint;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: <Widget>[
        Padding(
          padding: EdgeInsets.symmetric(horizontal: 5, vertical: 10),
          child: TextField(
            decoration: InputDecoration(
              border: OutlineInputBorder(),
              hintText: hint,
            ),
          ),
        ),
      ],
    );
  }
}

class CustomPasswdForm extends StatelessWidget {
  const CustomPasswdForm({
    super.key,
    required this.hint,
  });

  final String hint;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: <Widget>[
        Padding(
          padding: EdgeInsets.symmetric(horizontal: 5, vertical: 10),
          child: TextField(
            obscureText: true,
            decoration: InputDecoration(
              border: OutlineInputBorder(),
              hintText: hint,
            ),
          ),
        ),
      ],
    );
  }
}

class OrDivider extends StatelessWidget {
  const OrDivider({
    super.key,
    this.text = const AutoSizeText(''),
    this.textPadding = const EdgeInsets.symmetric(),
    this.color = Colors.grey,
    this.thickness = 1,
  });

  final AutoSizeText text;
  final EdgeInsets textPadding;
  final Color color;
  final double thickness;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Expanded(child: Divider(thickness: thickness, color: color,),),
        Padding(padding: textPadding, child: text,),
        Expanded(child: Divider(thickness: thickness, color: color,),),
      ],
    );
  }
}
